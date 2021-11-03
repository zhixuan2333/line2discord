/*
  Warnings:

  - A unique constraint covering the columns `[LineID]` on the table `Channel` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[DiscordID]` on the table `Channel` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "Channel_LineID_key" ON "Channel"("LineID");

-- CreateIndex
CREATE UNIQUE INDEX "Channel_DiscordID_key" ON "Channel"("DiscordID");
