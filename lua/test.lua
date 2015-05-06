synfunc = function (msg) msg:Reply("ack") end
robot:Respond("syn",synfunc)

--respond(robot.ID, "syn", synfunc)
